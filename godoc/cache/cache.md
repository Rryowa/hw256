TL;DR - Арк совмещает в себе LRU и MRU, а также учитывает возможность того, что нам пригодится кешированный  
запрос независимо от того сколько времени прошло с момента последнего такого запроса (условно - отгрузка  
заказов в пункт выдачи происходит утром, иногда может задерживаться, значит нам выгоднее инвалидировать кеш по популярности запросов)

TL;DR2 - Фактическое использование кэша в ARC ограничено размером кэша c. Этот адаптивный механизм гарантирует,  
что ARC динамически балансирует между recency(свежестью?{я не лингвист}) и частотой(frequency) данных, не превышая предварительно определенного размера кэша, 
что, по сути, ограничивает потребление ресурсов кэша.

Adaptive Replacement Cache (ARC), a novel cache management policy designed for demand paging systems.

The goal of ARC is to improve cache performance by dynamically balancing between recency and frequency in an online, self-tuning manner

ARC consistently outperforms LRU and performs comparably to more complex algorithms like LRU-2, 2Q, LRFU, and LIRS, even when these algorithms use the best offline tuning parameters:  
https://web.archive.org/web/20150405221102/https://www.usenix.org/legacy/event/fast03/tech/full_papers/megiddo/megiddo.pdf  

At any time, the behavior of the policy ARC is completely described once a certain adaptation parameter p [0<=p<=c] is known.
ARC continuously adapts and tunes p in response to an observed workload, which item to replace at any
given time.  
On a hit in B1, increment *p* by 1, if the size of B1 is at least the size of B2; otherwise, we increment *p* by module(B2)/module(B1).  
All increments to *p* are subject to a cap of *c*. Thus, smaller the size of B1, the larger the increment.  
Similarly, on a hit in B2, we decrement *p* by 1, if the size of B2 is at least the size of B1; otherwise, we decrement *p* module(B1)/module(B2).  
All decrements to *p* are subject to a minimum of 0. Thus, smaller the size of B2 , the larger the decrement.  
The idea is to “invest” in the list that is performing the best.